// /mypage/page.js
"use client"
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import CreateSession from './CreateSession';
import ErrorMessage from '@/components/ErrorMessage';
import { useErrorHandler } from "@/hooks/useErrorHandler";


export default function Mypage(){
    const [user, setUser] = useState(null);
    const router = useRouter();
    const [error,setError] = useState(null);
    const errorHandler = useErrorHandler(setError);

    //マイページ情報取得
    useEffect(() => {
    const load = async () => {
      try {
        const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/session/mypage`, 
        {
          credentials: "include",
        });

     if (!errorHandler(res)) return;

     const data = await res.json();
          setUser(data);

    } catch (e){
            setError("通信エラーが発生しました。");
            return ;
         }
    };

    load();
    },[]);

    //ログアウト処理
    const handleLogout = async() => {
      await fetch(`${process.env.NEXT_PUBLIC_API_URL}/auth/logout`,
        {
          method: "post",
        });
      router.push("/login");
    };

  if (error) return <ErrorMessage error={error} />;

    if(!user) return <p>読み込み中...</p>;

 return (
  <div className="container mt-4">
    <div className="d-flex justify-content-between align-items-center mb-4">
      <h1>{user.userName}さんのマイページ</h1>
      <button className="btn btn-danger" onClick={handleLogout}>ログアウト</button>
    </div>

    <CreateSession />

    <h3 className="mb-4"> 📊 テスト結果分析</h3>

    {user.testSessDtos.length === 0 ? (
      <div className="alert alert-info">
        テスト結果はまだありません。テストを開始すると結果が表示されます。
      </div>
    ) : (
      user.testSessDtos.map((session, index) => (
        <div key={index} className="mb-5 p-3 border rounded shadow-sm bg-secondary text-light">
          <h5>{session.startTime.split(' ')[0]}</h5>
          <p>正答数: {session.correctCount} / {session.total}</p>

          <h6>分野別正解数</h6>
          <div className="d-flex gap-3">
            {session.categoryDtos.map((cat, i) => (
                <div
                  key={i}
                  className="d-flex flex-column align-items-center"
                  style={{ width: '60px' }}
                >
                  <div
                    className="d-flex align-items-end"
                    style={{
                      height: '150px',
                      width: '100%',
                      backgroundColor: 'transparent',
                    }}
                  >
                    <div
                      className="bg-dark mx-auto"
                      style={{
                        height: `${(cat.correctCount / cat.total) * 100}%`,
                        width: '30px',
                        minHeight: '2px',
                        transition: 'height 0.3s ease-in-out',
                      }}
                    ></div>
                  </div>

                  <div className="text-center mt-2">
                    <small className="text-light d-block">
                      {cat.correctCount}/{cat.total}
                    </small>
                    <small>{cat.categoryName}</small>
                  </div>
                </div>
            ))}
          </div>

          <h6 className="mt-3">苦手分野（TOP2）</h6>
          <ul>
            {session.weakCategories.map((cat, i) => (
              <li key={i}>{cat}</li>
            ))}
          </ul>
        </div>
      ))
    )}
  </div>
);

}