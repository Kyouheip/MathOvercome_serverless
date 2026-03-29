import { fetchAuthSession } from "aws-amplify/auth";

// CognitoのIDトークンを取得してAPI Gatewayへの認証ヘッダーを返す
export async function getAuthHeader() {
  const session = await fetchAuthSession();
  const token = session.tokens?.idToken?.toString();
  return { Authorization: `Bearer ${token}` };
}
