import { useEffect, useState } from "react";

export default function useLocalStorage<T>(
  key: string,
  defaultValue: T
): [T, (value: T) => void] {
  const [value, setValue] = useState(defaultValue);

  useEffect(() => {
    const item = localStorage.getItem(key);

    if (!item) {
      localStorage.setItem(key, JSON.stringify(defaultValue));
    }

    setValue(item ? JSON.parse(item) : defaultValue);
  }, [key, defaultValue]);

  const setValueWrap = (value: T) => {
    try {
      setValue(value);
      localStorage.setItem(key, JSON.stringify(value));
    } catch (e) {
      console.error(e);
    }
  };

  return [value, setValueWrap];
}
