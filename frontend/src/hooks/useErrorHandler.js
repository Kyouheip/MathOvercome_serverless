// hooks/useErrorHandler.js
export function useErrorHandler(setError) {
  return (res) => {
    if (res.status === 401) {
      setError("セッションが切れたか不正なアクセスです。ログインし直してください。");
      return false;
    }
    if (!res.ok) {
      setError(`送信に失敗しました(${res.status})`);
      return false;
    }
    return true;
  };
}
