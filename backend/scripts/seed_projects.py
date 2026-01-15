import csv
import io
import json
import os
import sqlite3
import urllib.request
import ssl


DEFAULT_URLS = [
    (
        "https://docs.google.com/spreadsheets/d/"
        "1Md9w8s-K2qdA3HoGsxTeUtGg91QvtbdUrpYuk02dvto/"
        "export?format=csv&gid=1580934801"
    ),
    (
        "https://docs.google.com/spreadsheets/d/"
        "1Md9w8s-K2qdA3HoGsxTeUtGg91QvtbdUrpYuk02dvto/"
        "export?format=csv&gid=1815786844"
    ),
]


def normalize_username(value: str) -> str:
    value = value.strip()
    if value.startswith("@"):
        value = value[1:]
    return value.lower()


def is_project_header(row):
    header_role = row[4].strip()
    return (
        len(row) >= 8
        and row[0].strip()
        and row[2].strip() == "ФИО"
        and row[3].strip() == "Юзернейм"
        and header_role
        and header_role.startswith("Должность")
        and row[7].strip()
        and not row[0].strip().startswith("Руководство")
    )


def fetch_rows(url):
    context = ssl._create_unverified_context()
    with urllib.request.urlopen(url, context=context) as response:
        content = response.read().decode("utf-8")
    reader = csv.reader(io.StringIO(content))
    return list(reader)


def parse_projects(rows):
    projects = []
    current = None

    for row in rows:
        if len(row) < 8:
            row = row + [""] * (8 - len(row))

        if is_project_header(row):
            if current:
                projects.append(current)

            current = {
                "title": row[7].strip(),
                "status": row[0].strip(),
                "description": "",
                "members": [],
                "links": set(),
            }
            continue

        if not current:
            continue

        username = normalize_username(row[3])
        if username and username not in {"-", "—"}:
            member = {
                "username": username,
                "full_name": row[2].strip(),
                "role": row[4].strip(),
            }
            current["members"].append(member)

        link = row[7].strip()
        if link:
            current["links"].add(link)

    if current:
        projects.append(current)

    for project in projects:
        links = sorted(project.pop("links"))
        project["description"] = "\n".join(links)

    return projects


def seed_database(db_path, projects):
    os.makedirs(os.path.dirname(db_path), exist_ok=True)
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    cursor.execute(
        """
        CREATE TABLE IF NOT EXISTS projects (
            id INTEGER PRIMARY KEY,
            description TEXT,
            users TEXT,
            title TEXT,
            status TEXT
        );
        """
    )
    cursor.execute("DELETE FROM projects")

    for project in projects:
        cursor.execute(
            "INSERT INTO projects (description, users, title, status) VALUES (?, ?, ?, ?)",
            (
                project["description"],
                json.dumps(project["members"], ensure_ascii=False),
                project["title"],
                project["status"],
            ),
        )

    conn.commit()
    conn.close()


def main():
    url_env = os.getenv("PROJECTS_SPREADSHEET_URLS", "")
    if url_env:
        urls = [value.strip() for value in url_env.split(",") if value.strip()]
    else:
        urls = DEFAULT_URLS

    db_path = os.getenv(
        "PROJECTS_DB_PATH",
        os.path.abspath(
            os.path.join(
                os.path.dirname(__file__),
                "..",
                "internal",
                "model",
                "database",
                "projects_db.db",
            )
        ),
    )

    projects = []
    for url in urls:
        rows = fetch_rows(url)
        projects.extend(parse_projects(rows))

    seed_database(db_path, projects)
    print(f"Seeded {len(projects)} projects into {db_path}")


if __name__ == "__main__":
    main()
