//problems/QuestionForm.js
'use client'
import {useEffect, useState} from 'react'
import { useRouter } from 'next/navigation'
import { useErrorHandler } from "@/hooks/useErrorHandler"
import ErrorMessage from '@/components/ErrorMessage';

export default function QuestionForm({idx,choices,initialselectedId,total}){
    const [selectedId,setSelectedId] = useState(initialselectedId ?? null);
    const router = useRouter();
    const [error,setError] = useState(null);
    const errorHandler = useErrorHandler(setError);

     useEffect(() => {
    setSelectedId(initialselectedId ?? null);
            }, [idx, initialselectedId]);

    const submit = async (e) => {
        //formを使うときは必要。ないとリロードされReactが無視される
        e.preventDefault();

        try{
            //選択肢を送る
        const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/session/current/problems/${idx}/answer`,
            {
                method: "post",
                headers: {"Content-Type": "application/json"},
                credentials: 'include',
                body: JSON.stringify({selectedChoiceId: selectedId}),
            }
        );

       if (!errorHandler(res)) return;

        }catch (e) {
             setError("通信エラーが発生しました");
             return ;
        }
        const nextIdx = idx+1;
        if(nextIdx < total){
            router.push(`/problems?idx=${nextIdx}`)
        } else{
            router.push(`/mypage`)
        }
    };

    const handleBack = () => {
        router.push(`/problems?idx=${Math.max(0, idx - 1)}`);
    }

    if (error) return <ErrorMessage error={error} />;

    return(
        <form onSubmit={submit}>
          <div className="mb-3">
            {choices.map(choice => (
                <div className="form-check" key={choice.id}>
                    <input
                     className="form-check-input"
                     type="radio"
                     name="choice"
                     id={`choice-${choice.id}`}
                     checked={selectedId === choice.id}
                     onChange={() => setSelectedId(choice.id)}
                     />
                    <label htmlFor={`choice-${choice.id}`} className="form-check-label">
                      {choice.choiceText}
                    </label>
                </div>
                    )
                )
            }
           </div>

            <div className="mt-4">
                {Number(idx) > 0 && (
                    <button type="button" className="btn btn-secondary me-2" onClick={handleBack}>戻る</button>
                )}
                <button type="submit" className="btn btn-primary">次へ</button>
            </div>
         </form>
    );
}
