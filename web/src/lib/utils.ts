import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatAmount(amount: number, type: 'income' | 'expense'): string {
  const prefix = type === 'income' ? '+' : '-';
  return `${prefix}${amount.toFixed(2)}`;
}

/**
 * Safe date formatter that returns a fallback string instead of throwing
 * on invalid dates (prevents RangeError: Invalid time value crashes).
 */
export function safeFormat(
  dateValue: string | Date | null | undefined,
  formatter: (d: Date) => string,
  fallback = '—'
): string {
  if (!dateValue) return fallback;
  try {
    const d = dateValue instanceof Date ? dateValue : new Date(dateValue);
    if (isNaN(d.getTime())) return fallback;
    return formatter(d);
  } catch {
    return fallback;
  }
}

/**
 * Safe toFixed that coerces string/unknown values to number first.
 * API may return amounts as strings (e.g. "100.00") instead of numbers.
 */
export function safeFixed(value: unknown, digits = 2): string {
  const n = Number(value);
  if (isNaN(n)) return '0.' + '0'.repeat(digits);
  return n.toFixed(digits);
}
