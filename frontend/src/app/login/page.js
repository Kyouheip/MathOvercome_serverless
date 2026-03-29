// /login/page.js
"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { signIn, signOut } from "aws-amplify/auth";

export default function LoginPage() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const router = useRouter();

  const doSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await signOut().catch(() => {});
      await signIn({ username: email, password });
      router.push("/mypage");
    } catch (err) {
      setError(err.message || "ログインに失敗しました");
    }
  };

  return (
    <div className="container mt-4">
      <h2>ログイン</h2>
      <form onSubmit={doSubmit}>
        <div className="mb-3">
          <label htmlFor="email" className="form-label">メールアドレス</label>
          <input
            id="email"
            type="email"
            className="form-control"
            value={email}
            onChange={e => setEmail(e.target.value)}
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
        <p className="mt-2">新規登録は<a href="/register">こちら</a></p>
      </form>
    </div>
  );
}
