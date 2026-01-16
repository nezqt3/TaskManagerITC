export function getJwt() {
  return localStorage.getItem("jwt");
}

export function getProfile() {
  const stored = localStorage.getItem("profile");
  if (!stored) {
    return null;
  }

  try {
    return JSON.parse(stored);
  } catch (error) {
    return null;
  }
}

export function getAuthHeaders() {
  const jwt = getJwt();
  if (!jwt) {
    return {};
  }
  return { Authorization: `Bearer ${jwt}` };
}

export function parseRoles(role = "") {
  return role
    .toLowerCase()
    .split(/[\\s,;/|+]+/)
    .map((item) => item.trim())
    .filter(Boolean);
}

export function isAdmin(role = "") {
  const roles = parseRoles(role);
  return roles.includes("админ") || roles.includes("admin") || roles.includes("владелец");
}

export function isModerator(role = "") {
  const roles = parseRoles(role);
  return roles.includes("модератор") || roles.includes("moderator");
}

export function isLeader(role = "") {
  const roles = parseRoles(role);
  return roles.includes("руководитель");
}

export function canManageMembers(role = "") {
  return isAdmin(role);
}

export function canManageTasks(role = "") {
  return isAdmin(role) || isModerator(role) || isLeader(role);
}

export function canReview(role = "") {
  return canManageTasks(role);
}
