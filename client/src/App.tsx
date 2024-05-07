import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

    const handleLogin = () => {
      window.location.href = "http://localhost:8080/auth/google";
    }

    const handleLogout = async () => {
        try {
            const response = await fetch("http://localhost:8080/auth/google/logout", {
                method: 'GET', // Assuming logout requires a POST request
            });

            if (!response.ok) {
                // Handle errors here, e.g., display an error message to the user
                console.error("Logout request failed:", response.statusText);
                return;
            }

            // Handle successful logout response (optional)
            console.log("Logout successful!");

        } catch (error) {
            console.error("Error during logout request:", error);
        }
    };

  return (
      <>
          <div>
              <a href="https://vitejs.dev" target="_blank">
                  <img src={viteLogo} className="logo" alt="Vite logo"/>
              </a>
              <a href="https://react.dev" target="_blank">
                  <img src={reactLogo} className="logo react" alt="React logo"/>
              </a>
          </div>
          <h1>Vite + React</h1>
          <div className="card">
              <button onClick={() => setCount((count) => count + 1)}>
                  count is {count}
              </button>
              <p>
                  Edit <code>src/App.tsx</code> and save to test HMR
              </p>
          </div>

          <div>
              <button onClick={handleLogin}>Login with Google</button>
          </div>
          <div>
              <button onClick={handleLogout}>Logout from Google</button>
          </div>

          <p className="read-the-docs">
              Click on the Vite and React logos to learn more
          </p>
      </>
  )
}

export default App
