// /crs/app
import "./globals.css";

export const metadata = {
  title: "MathOvercome"
};

export default function RootLayout({ children }) {
  return (
    <html lang="ja">
      <head>
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/bootstrap.min.css"
                    rel="stylesheet"
                    integrity="sha384-9ndCyUaIbzAi2FUVXJi0CjmCapSmO7SnpJef0486qhLnuZ2cdeRhO02iuK6FUUVM"
                    crossOrigin="anonymous"/>
      </head>

      <body className="bg-dark text-white">
        <header className="bg-secondary text-dark p-3 position-relative">
          <div className="d-flex flex-column flex-md-row justify-content-between align-items-md-end">
            <h1 className="mb-2 mb-md-0">MathOvercome</h1>
             <div className=" text-dark fs-6">
              ～数学の苦手を克服し大学受験を乗り越えよう！
             </div>
          </div>
        </header>

        <main>{children}</main>
        
      </body>
    </html>
  );
}
