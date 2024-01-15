export const ShortDate = (date) =>
  new Date(date).toLocaleDateString("en-UK", {
    day: "numeric",
    month: "short",
    year: "2-digit",
  });

export const LongDate = (date) =>
  new Date(date).toLocaleDateString("en-UK", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });

export const BirthdayConverter = (date) => {
  if (!date) {
    return;
  }
  const [day, month, year] = date?.split("/");
  return LongDate(year, month - 1, day);
};

export const ShortDatetime = (datetime) => (
  <span className="text-nowrap">
    {new Date(datetime).toLocaleTimeString("en-UK", {
      year: "2-digit",
      month: "short",
      day: "2-digit",
      hour: "numeric",
      minute: "2-digit",
    })}
  </span>
);

export const ShortTime = (time) =>
  new Date(time).toLocaleTimeString([], {
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  });
