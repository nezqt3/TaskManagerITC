export const normalizeUsername = (value = "") =>
  value.toString().trim().replace(/^@/, "").toLowerCase();

export const normalizeFullName = (value = "") =>
  value.toString().trim().replace(/\s+/g, " ").toLowerCase();

export const formatFullName = ({ fullName, firstName, lastName }) => {
  const parts = fullName ? fullName.trim().split(/\s+/).filter(Boolean) : [];
  const first = firstName || parts[0] || "";
  const last = lastName || parts[1] || "";
  const combined = [first, last].filter(Boolean).join(" ");

  if (combined) {
    return combined;
  }
  if (parts.length > 0) {
    return parts.join(" ");
  }
  return "";
};

export const formatDisplayName = ({ fullName, firstName, lastName, username }) => {
  const base = formatFullName({ fullName, firstName, lastName });
  const nick = normalizeUsername(username || "");

  if (base && nick) {
    return `${base} (${nick})`;
  }
  if (base) {
    return base;
  }
  if (nick) {
    return nick;
  }
  return "";
};
