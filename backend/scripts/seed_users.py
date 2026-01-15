import csv
import io
import os
import sqlite3


def normalize_username(value: str) -> str:
    value = value.strip()
    if value.startswith("@"):
        value = value[1:]
    return value


def fetch_csv(path: str) -> list[list[str]]:
    with open(path, "r", encoding="utf-8") as handle:
        content = handle.read()
    return list(csv.reader(io.StringIO(content)))


def parse_users(rows: list[list[str]]) -> list[dict]:
    users: list[dict] = []
    for i, row in enumerate(rows):
        if i in (0, 1):
            continue
        if len(row) < 7:
            row = row + [""] * (7 - len(row))
        telegram_id = row[5].strip()
        if not telegram_id:
            continue
        users.append(
            {
                "full_name": row[0].strip(),
                "username": normalize_username(row[1]),
                "date_of_birthday": row[2].strip(),
                "number_of_phone": row[3].strip(),
                "role": row[4].strip(),
                "telegram_id": telegram_id,
                "may_to_open": row[6].strip().lower() == "true",
            }
        )
    return users


def seed_database(db_path: str, users: list[dict]) -> None:
    os.makedirs(os.path.dirname(db_path), exist_ok=True)
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()

    cursor.execute(
        """
        CREATE TABLE IF NOT EXISTS users (
            TelegramID TEXT PRIMARY KEY,
            FirstName TEXT,
            LastName TEXT,
            Username TEXT,
            PhotoURL TEXT,
            FullName TEXT,
            DateOfBirthday TEXT,
            NumberOfPhone TEXT,
            Role TEXT,
            MayToOpen BOOLEAN
        );
        """
    )
    cursor.execute("DELETE FROM users")

    for user in users:
        cursor.execute(
            """
            INSERT INTO users (
                TelegramID,
                FirstName,
                LastName,
                Username,
                PhotoURL,
                FullName,
                DateOfBirthday,
                NumberOfPhone,
                Role,
                MayToOpen
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
            """,
            (
                user["telegram_id"],
                "",
                "",
                user["username"],
                "",
                user["full_name"],
                user["date_of_birthday"],
                user["number_of_phone"],
                user["role"],
                user["may_to_open"],
            ),
        )

    conn.commit()
    conn.close()


def main() -> None:
    source = os.getenv("USERS_CSV_PATH", "").strip()
    if not source:
        raise SystemExit("USERS_CSV_PATH is required")

    db_path = os.getenv(
        "USERS_DB_PATH",
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

    rows = fetch_csv(source)
    users = parse_users(rows)
    seed_database(db_path, users)
    print(f"Seeded {len(users)} users into {db_path}")


if __name__ == "__main__":
    main()
