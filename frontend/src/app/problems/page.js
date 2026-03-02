// problems/page.js
"use client";
import { Suspense, useEffect, useState } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import QuestionForm from "./QuestionForm";
import ErrorMessage from "@/components/ErrorMessage";
import { useErrorHandler } from "@/hooks/useErrorHandler";

function ProblemsContent() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const idx = Number(searchParams.get("idx") ?? 0);

  const [sp, setSp] = useState(null);
  const [error, setError] = useState(null);
  const [showHint, setShowHint] = useState(false);
  const errorHandler = useErrorHandler(setError);

  useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/session/current/problems/${idx}`,
          { credentials: "include" }
        );

        if (!errorHandler(res)) return;

        const data = await res.json();
        setSp(data);
      } catch {
        setError("通信エラーが発生しました。");
      }
    };
    load();
  }, [idx]);

  useEffect(() => {
    setShowHint(false);
  }, [idx]);


  if (error) return <ErrorMessage error={error} />;
  if (!sp) return <p>読み込み中...</p>;

  return (
    <div className="container mt-4">
      <h2 className="mb-3">問題 {idx + 1} / {sp.total}</h2>

      <div className="bg-secondary p-3 rounded mb-4">
        <p
          className="mb-0 text-dark fw-bold"
          dangerouslySetInnerHTML={{ __html: sp.question }}
        />
      </div>

      <QuestionForm
        idx={idx}
        choices={sp.choices}
        initialselectedId={sp.selectedId}
        total={sp.total}
        onNext={() => router.push(`/problems?idx=${idx + 1}`)}
        onPrev={() => router.push(`/problems?idx=${Math.max(0, idx - 1)}`)}
      />

      <div className="mt-4">
        <button className="btn btn-outline-info" onClick={() => setShowHint(!showHint)}>
          ヒント！
        </button>

        {showHint && (
          <div className="mt-3 p-3 rounded bg-secondary text-dark">
            <strong>ヒント：</strong> {sp.hint}
          </div>
        )}
      </div>
    </div>
  );
}

export default function ProblemsPage() {
  return (
    <Suspense fallback={<p>読み込み中...</p>}>
      <ProblemsContent />
    </Suspense>
  );
}
