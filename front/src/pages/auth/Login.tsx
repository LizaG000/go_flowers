import { useState } from "react"
import { v4 as uuidv4 } from "uuid"
import { useNavigate } from "react-router-dom"
import "../../styles/Registration.css"

function Login() {
  const [inputName, setInputName] = useState("")
  const [inputPassword, setInputPassword] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState("")
  const navigate = useNavigate()

  const login = async () => {
    const username = inputName.trim()
    const password = inputPassword.trim()

    setErrorMessage("")

    if (!username || !password || isLoading) {
      setErrorMessage("Заполните все поля")
      return
    }

    setIsLoading(true)

    try {
      const response = await fetch("/api/auth/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Idempotency-Key": uuidv4(),
        },
        body: JSON.stringify({
          username,
          password,
        }),
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(JSON.stringify(data))
      }

      localStorage.setItem("access_token", data.access_token)
      localStorage.setItem("token_type", data.token_type)

      if (data.user) {
        localStorage.setItem("user", JSON.stringify(data.user))
      }

      navigate("/")
    } catch (error) {
      console.error(error)
      setErrorMessage(
        "Не удалось авторизоваться. Проверьте введённые данные.",
      )
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <main className="registration-page">
      <div className="registration-card">
        <h1>Вход</h1>

        <p className="registration-card__subtitle">
          Введите данные для входа
        </p>

        <div className="registration-field">
          <label htmlFor="name">Username</label>

          <input
            type="text"
            id="name"
            value={inputName}
            onChange={(event) => setInputName(event.target.value)}
            required
          />
        </div>

        <div className="registration-field">
          <label htmlFor="password">Password</label>

          <input
            type="password"
            id="password"
            value={inputPassword}
            onChange={(event) => setInputPassword(event.target.value)}
            required
          />
        </div>

        {errorMessage && (
          <p className="registration-message registration-message--error">
            {errorMessage}
          </p>
        )}

        <button
          className="registration-button"
          type="button"
          onClick={login}
          disabled={isLoading}
        >
          {isLoading ? "Авторизация..." : "Войти"}
        </button>
      </div>
    </main>
  )
}

export default Login