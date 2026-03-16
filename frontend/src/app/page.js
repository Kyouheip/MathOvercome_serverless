// src/app/page.js
'use client'
import { useRouter } from 'next/navigation'

export default function HomePage() {
  const router = useRouter();

  const handleLogin = () => {
    router.push('/login');
  }

  return (
    <main className="container mt-5 text-center">
      <h1 className="mb-4">高校数学の苦手を分析して苦手をなくそう！</h1>
      <p className="lead">共通試験の標準レベルの問題を基準にしています</p>
      <p className="mb-4">さあ始めよう！</p>
      <button className="btn btn-primary btn-lg" onClick={handleLogin}>
        ログイン画面へ
      </button>
    </main>
  )
}
