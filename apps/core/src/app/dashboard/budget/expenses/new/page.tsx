"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

export default function NewExpensePage() {
  const router = useRouter();

  useEffect(() => {
    // redirect to budget dashboard with query param to open expense dialog
    router.replace("/dashboard/budget?openExpense=true");
  }, [router]);

  // return null while redirecting
  return null;
}
