const DEADLINE_EMPTY_LABEL = "без ограничений";

const getIsoDateParts = (value) => {
  if (typeof value !== "string") {
    return null;
  }
  const match = value.trim().match(/^(\d{4})-(\d{2})-(\d{2})/);
  if (!match) {
    return null;
  }
  return {
    year: Number(match[1]),
    yearRaw: match[1],
    month: match[2],
    day: match[3],
  };
};

const getDmyDateParts = (value) => {
  if (typeof value !== "string") {
    return null;
  }
  const match = value.trim().match(/^(\d{2})\.(\d{2})\.(\d{4})$/);
  if (!match) {
    return null;
  }
  return {
    year: Number(match[3]),
    yearRaw: match[3],
    month: match[2],
    day: match[1],
  };
};

const isEmptyDeadline = (value) => {
  if (!value) {
    return true;
  }
  const parts = getIsoDateParts(value);
  if (parts && parts.year <= 1) {
    return true;
  }
  const dmyParts = getDmyDateParts(value);
  if (dmyParts && dmyParts.year <= 1) {
    return true;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return true;
  }
  return date.getUTCFullYear() <= 1;
};

export const getDeadlineDate = (value) => {
  if (isEmptyDeadline(value)) {
    return null;
  }
  const dmyParts = getDmyDateParts(value);
  if (dmyParts) {
    const date = new Date(
      `${dmyParts.yearRaw}-${dmyParts.month}-${dmyParts.day}T00:00:00`
    );
    return Number.isNaN(date.getTime()) ? null : date;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return null;
  }
  return date;
};

export const formatDeadline = (value) => {
  if (isEmptyDeadline(value)) {
    return DEADLINE_EMPTY_LABEL;
  }
  const dmyParts = getDmyDateParts(value);
  if (dmyParts) {
    return `${dmyParts.day}.${dmyParts.month}.${dmyParts.yearRaw}`;
  }
  const parts = getIsoDateParts(value);
  if (parts) {
    return `${parts.day}.${parts.month}.${parts.yearRaw}`;
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return DEADLINE_EMPTY_LABEL;
  }
  const day = String(date.getDate()).padStart(2, "0");
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const year = String(date.getFullYear()).padStart(4, "0");
  return `${day}.${month}.${year}`;
};
