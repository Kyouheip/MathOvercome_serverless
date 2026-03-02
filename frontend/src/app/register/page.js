// /register/page.js
"use client"
import {useState} from "react";
import {useRouter} from "next/navigation"

export default function RegisterPage(){
    const[userName,setUserName] = useState("");
    const [userId,setUserId] = useState("");
    const [password1,setPassword1] = useState("");
    const[password2,setPassword2] = useState("");
    const[error,setError] = useState("");
    const router = useRouter();

    const doRegister = async (e) => {
        e.preventDefault();
        setError("");

        const res = await fetch(
            `${process.env.NEXT_PUBLIC_API_URL}/auth/register`,
            {
                method: "POST",
                headers: {"Content-Type":"application/json"},
                body: JSON.stringify({userName,userId,password1,password2}),
            }
        );

        if(res.status === 201){
            router.push("/login");
        }else{
            const msg = await res.text();
            setError(msg || `エラー: ${res.status}`);
        }
    };

    const handleBack = () => {
        router.push("/login");
    }

    return(
        <div className="container mt-4">
            <h2 className="mb-3">新規登録</h2>
            <form onSubmit = {doRegister}>
            <div className="mb-3">
              <label htmlFor="userId" className="form-label">ID</label>
              <input
                id="userId" //css
                type="text"
                className="form-control" //css
                value={userId}
                onChange={e => setUserId(e.target.value)}
              />
              </div>

              <div className="mb-3">
                <label htmlFor="userName" className="form-label">名前</label>
                <input
                    id="userName"
                    type="text"
                    className="form-control"
                    value={userName}
                    onChange={e => setUserName(e.target.value)}
                />
              </div>
            
              <div className="mb-3">
                <label htmlFor="password1" className="form-label">パスワード</label>
                <input
                    id="password1"
                    type="password"
                    className="form-control"
                    value={password1}
                    onChange={e => setPassword1(e.target.value)}
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
            <button type="button" className="btn btn-secondary" onClick={handleBack}>戻る</button>
            </form>
        </div>
    );

}
