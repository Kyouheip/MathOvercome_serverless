// /register/page.js
"use client";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { signUp, confirmSignUp } from "aws-amplify/auth";

export default function RegisterPage() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [password2, setPassword2] = useState("");
  const [confirmCode, setConfirmCode] = useState("");
  const [step, setStep] = useState("register"); // "register" | "confirm"
  const [error, setError] = useState("");
  const router = useRouter();

  const doRegister = async (e) => {
    e.preventDefault();
    setError("");

    if (password !== password2) {
      setError("パスワードが一致しません");
      return;
    }

    try {
      await signUp({
        username: email,
        password,
        options: {
          userAttributes: {
            email,
            name,
          },
        },
      });
      setStep("confirm");
    } catch (err) {
      setError(err.message || "登録に失敗しました");
    }
  };

  const doConfirm = async (e) => {
    e.preventDefault();
    setError("");
    try {
      await confirmSignUp({ username: email, confirmationCode: confirmCode });
      router.push("/login");
    } catch (err) {
      setError(err.message || "確認に失敗しました");
    }
  };

  if (step === "confirm") {
    return (
      <div className="container mt-4">
        <h2 className="mb-3">メール確認</h2>
        <p>{email} に確認コードを送信しました。</p>
        <form onSubmit={doConfirm}>
          <div className="mb-3">
            <label htmlFor="confirmCode" className="form-label">確認コード</label>
            <input
              id="confirmCode"
              type="text"
              className="form-control"
              value={confirmCode}
              onChange={e => setConfirmCode(e.target.value)}
            />
          </div>
          {error && <pre className="text-danger">{error}</pre>}
          <button type="submit" className="btn btn-primary">確認</button>
        </form>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <h2 className="mb-3">新規登録</h2>
      <form onSubmit={doRegister}>
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
          <label htmlFor="name" className="form-label">名前</label>
          <input
            id="name"
            type="text"
            className="form-control"
            value={name}
            onChange={e => setName(e.target.value)}
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

        <div className="mb-3">
          <label htmlFor="password2" className="form-label">パスワード確認</label>
          <input
            id="password2"
            type="password"
            className="form-control"
            value={password2}
            onChange={e => setPassword2(e.target.value)}
          />
        </div>

        {error && <pre className="text-danger">{error}</pre>}
        <button type="submit" className="btn btn-primary me-2">登録</button>
        <button type="button" className="btn btn-secondary" onClick={() => router.push("/login")}>戻る</button>
      </form>
    </div>
  );
}
