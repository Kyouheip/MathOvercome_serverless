// components/ErrorMessage.js
export default function ErrorMessage({ error }) {
  if (!error) return null;
  return (
    <div>
      <p>{error}</p>
      <p><a href="/login">ログインページへ</a></p>
    </div>
  );
}