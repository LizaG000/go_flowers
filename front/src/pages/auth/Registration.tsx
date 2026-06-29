import { useState } from "react"
import { v4 as uuidv4 } from "uuid"
import { useNavigate } from "react-router-dom"
import "../../styles/Registration.css"

function Registration() {
    const [inputName, setInputName] = useState("")
    const [inputEmail, setInputEmail] = useState("")
    const [inputPassword, setInputPassword] = useState("")
    const [isLoading, setIsLoading] = useState(false)
    const [errorMessage, setErrorMessage] = useState("")
    const [successMessage, setSuccessMessage] = useState("")
    const navigate = useNavigate()



    const register = async () => {
        const username = inputName.trim()
        const email = inputEmail.trim()
        const password = inputPassword.trim()

        setErrorMessage("")
        setSuccessMessage("")

        if (!username || !email || !password || isLoading) {
            setErrorMessage("Заполните все поля")
            return
        }

        setIsLoading(true)

        try {
            const response = await fetch("/api/auth/register", {
                method: "POST",
                body: JSON.stringify({
                    username,
                    email,
                    password,
                }),
                headers: {
                    "Content-Type": "application/json",
                    "Idempotency-Key": uuidv4(),
                },
            })

            const data = await response.json()

            if (!response.ok) {
                throw new Error(JSON.stringify(data))
            }

            localStorage.setItem("access_token", data.access_token)
            navigate("/")
        } catch (error) {
            console.error(error)

            setErrorMessage(
                "Не удалось зарегистрироваться. Проверьте введённые данные.",
            )
        } finally {
            setIsLoading(false)
        }
    }

    return (
  <main className="registration-page">
    <div className="registration-card">
      <h1>Регистрация</h1>
      <p className="registration-card__subtitle">
        Создайте аккаунт для работы с музыкой
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
        <label htmlFor="email">Email</label>
        <input
          type="email"
          id="email"
          value={inputEmail}
          onChange={(event) => setInputEmail(event.target.value)}
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

      {errorMessage && <p className="registration-message registration-message--error">{errorMessage}</p>}

      {successMessage && <p className="registration-message registration-message--success">{successMessage}</p>}

      <button
        className="registration-button"
        type="button"
        onClick={register}
        disabled={isLoading}
      >
        {isLoading ? "Регистрация..." : "Зарегистрироваться"}
      </button>
    </div>
  </main>
)
}

export default Registration