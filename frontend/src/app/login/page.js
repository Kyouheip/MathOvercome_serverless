// /login/page.js
"use client";
import {useState} from "react";
import {useRouter} from "next/navigation"

export default function LoginPage(){
  const [userId,setUserId] = useState("");
  const [password,setPassword] = useState("");
  const [error,setError] = useState("");
  const [error2,setError2] = useState("");
  const router = useRouter();

  const doSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try{
    const res = await fetch(
      `${process.env.NEXT_PUBLIC_API_URL}/auth/login`,
      {
      method: "POST",
      headers: {"Content-Type":"application/json"},
      credentials: "include",//セッションidを保存したCookieを受け取る
      body: JSON.stringify({userId,password}),
      }
    );

    if(res.ok){
      //ログイン成功。マイページへ
      router.push("/mypage");
    }else{
      const msg = await res.text();
      setError(msg || `エラー: ${res.status}`);
    }
   }catch (e) {
    setError2("通信エラーが発生しました");
        return ;
   }
  }


  if (error2) {
    return <p><a href="" onClick={() => window.location.reload()}>再読み込み</a></p>
  }

  return(
    <div className="container mt-4">
      <h2>ログイン</h2>
      <form onSubmit = {doSubmit}>
      <div className="mb-3">
        <label htmlFor="userId" className="form-label">ID</label>
        <input 
          id="userId"
          type = "text"
          className="form-control"
          value = {userId}
          onChange = {e => setUserId(e.target.value)}
        />
      </div>
      <div className="mb-3">
        <label htmlFor="password" className="form-label">パスワード</label>
        <input 
          id="password"
          type="password" 
          className="form-control"
          value={password} 
          onChange={e => setPassword(e.target.value)}
        />
      </div>
        {error && <pre className="text-danger">{error}</pre>}
        <button type="submit" className="btn btn-primary">ログイン</button>
      <p className="mt-2">新規登録は
      <a href="/register">こちら</a>
      </p>
      </form>
      </div>
  );
}
