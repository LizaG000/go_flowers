import { useEffect, useState } from "react"
import { v4 as uuidv4 } from "uuid"
import "../styles/Main.css"

type Song = {
  id: number
  name: string
  artist: string
  genre: string
  duration: number | null
}

const LIMIT = 10

function Main() {
  const [songs, setSongs] = useState<Song[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState("")
  const [page, setPage] = useState(1)
  const [pageMessage, setPageMessage] = useState("")

  // Поля создания песни
  const [newId, setNewId] = useState("")
  const [newName, setNewName] = useState("")
  const [newArtist, setNewArtist] = useState("")
  const [newGenre, setNewGenre] = useState("")
  const [newDuration, setNewDuration] = useState("")
  const [isCreating, setIsCreating] = useState(false)
  const [createMessage, setCreateMessage] = useState("")

  // Поля редактирования
  const [editingId, setEditingId] = useState<number | null>(null)
  const [editName, setEditName] = useState("")
  const [editArtist, setEditArtist] = useState("")
  const [editGenre, setEditGenre] = useState("")
  const [editDuration, setEditDuration] = useState("")
  const [actionMessage, setActionMessage] = useState("")

  const formatDuration = (seconds: number) => {
    const minutes = Math.floor(seconds / 60)
    const remainingSeconds = seconds % 60

    return `${minutes}:${String(remainingSeconds).padStart(2, "0")}`
  }

  const fetchSongsPage = async (pageNumber: number) => {
    const token = localStorage.getItem("access_token")
    const offset = (pageNumber - 1) * LIMIT

    const response = await fetch(
      `/api/songs/?offset=${offset}&limit=${LIMIT}`,
      {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    )

    if (!response.ok) {
      throw new Error("Не удалось получить список песен")
    }

    const data: Song[] = await response.json()
    return data
  }

  useEffect(() => {
    const loadSongs = async () => {
      setIsLoading(true)
      setError("")

      try {
        const data = await fetchSongsPage(page)
        setSongs(data)
      } catch (error) {
        console.error(error)
        setError("Не удалось загрузить песни")
      } finally {
        setIsLoading(false)
      }
    }

    loadSongs()
  }, [page])

  const createSong = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    const id = Number(newId)
    const duration = Number(newDuration)

    setCreateMessage("")

    if (!newId || !newName.trim() || !newArtist.trim() || !newGenre.trim()) {
      setCreateMessage("Заполните все обязательные поля")
      return
    }

    if (!Number.isInteger(id) || id <= 0) {
      setCreateMessage("ID должен быть целым числом больше 0")
      return
    }

    if (newDuration && (!Number.isInteger(duration) || duration < 0)) {
      setCreateMessage("Длительность должна быть числом в секундах")
      return
    }

    setIsCreating(true)

    try {
      const token = localStorage.getItem("access_token")

      const response = await fetch("/api/songs/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "Idempotency-Key": uuidv4(),
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          id,
          name: newName.trim(),
          artist: newArtist.trim(),
          genre: newGenre.trim(),
          duration: newDuration ? duration : null,
        }),
      })

      const data: Song = await response.json()

      if (!response.ok) {
        throw new Error("Не удалось добавить песню")
      }

      setCreateMessage("Песня успешно добавлена")

      setNewId("")
      setNewName("")
      setNewArtist("")
      setNewGenre("")
      setNewDuration("")

      if (page === 1) {
        setSongs((currentSongs) => [data, ...currentSongs].slice(0, LIMIT))
      } else {
        setPage(1)
      }
    } catch (error) {
      console.error(error)
      setCreateMessage("Не удалось добавить песню")
    } finally {
      setIsCreating(false)
    }
  }

  const startEditing = (song: Song) => {
    setEditingId(song.id)
    setEditName(song.name)
    setEditArtist(song.artist)
    setEditGenre(song.genre)
    setEditDuration(song.duration?.toString() ?? "")
    setActionMessage("")
  }

  const cancelEditing = () => {
    setEditingId(null)
    setEditName("")
    setEditArtist("")
    setEditGenre("")
    setEditDuration("")
  }

  const updateSong = async (songId: number) => {
    const duration = Number(editDuration)

    if (!editName.trim() || !editArtist.trim() || !editGenre.trim()) {
      setActionMessage("Заполните название, исполнителя и жанр")
      return
    }

    if (editDuration && (!Number.isInteger(duration) || duration < 0)) {
      setActionMessage("Длительность должна быть числом в секундах")
      return
    }

    try {
      const token = localStorage.getItem("access_token")

      const response = await fetch(`/api/songs/${songId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          id: songId,
          name: editName.trim(),
          artist: editArtist.trim(),
          genre: editGenre.trim(),
          duration: editDuration ? duration : null,
        }),
      })

      const updatedSong: Song = await response.json()

      if (!response.ok) {
        throw new Error("Не удалось изменить песню")
      }

      setSongs((currentSongs) =>
        currentSongs.map((song) =>
          song.id === songId ? updatedSong : song,
        ),
      )

      setActionMessage("Песня успешно изменена")
      cancelEditing()
    } catch (error) {
      console.error(error)
      setActionMessage("Не удалось изменить песню")
    }
  }

  const deleteSong = async (songId: number) => {
    const confirmed = window.confirm("Удалить эту песню?")

    if (!confirmed) {
      return
    }

    try {
      const token = localStorage.getItem("access_token")

      const response = await fetch(`/api/songs/${songId}`, {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error("Не удалось удалить песню")
      }

      const updatedSongs = songs.filter((song) => song.id !== songId)
      setSongs(updatedSongs)

      setActionMessage("Песня удалена")

      // Если после удаления страница стала пустой — возвращаемся назад
      if (updatedSongs.length === 0 && page > 1) {
        setPage(page - 1)
      }
    } catch (error) {
      console.error(error)
      setActionMessage("Не удалось удалить песню")
    }
  }

  const goToPreviousPage = () => {
    if (page > 1) {
      setPageMessage("")
      setPage(page - 1)
    }
  }

  const goToNextPage = async () => {
    const nextPage = page + 1

    setPageMessage("")

    try {
      const data = await fetchSongsPage(nextPage)

      // Следующая страница пустая:
      // остаёмся на текущей странице
      if (data.length === 0) {
        setPageMessage("Дальше песен нет")
        return
      }

      setSongs(data)
      setPage(nextPage)
    } catch (error) {
      console.error(error)
      setPageMessage("Не удалось загрузить следующую страницу")
    }
  }

  return (
    <main className="songs-page">
      <h1>Песни</h1>

      <form className="song-form" onSubmit={createSong}>
        <h2>Добавить песню</h2>

        <input
          type="number"
          min="1"
          placeholder="ID"
          value={newId}
          onChange={(event) => setNewId(event.target.value)}
          required
        />

        <input
          type="text"
          placeholder="Название"
          value={newName}
          onChange={(event) => setNewName(event.target.value)}
          required
        />

        <input
          type="text"
          placeholder="Исполнитель"
          value={newArtist}
          onChange={(event) => setNewArtist(event.target.value)}
          required
        />

        <input
          type="text"
          placeholder="Жанр"
          value={newGenre}
          onChange={(event) => setNewGenre(event.target.value)}
          required
        />

        <input
          type="number"
          min="0"
          placeholder="Длительность в секундах"
          value={newDuration}
          onChange={(event) => setNewDuration(event.target.value)}
        />

        <button type="submit" disabled={isCreating}>
          {isCreating ? "Добавление..." : "Добавить песню"}
        </button>

        {createMessage && <p>{createMessage}</p>}
      </form>

      {actionMessage && <p>{actionMessage}</p>}

      {isLoading && <p>Загрузка песен...</p>}

      {error && <p>{error}</p>}

      {!isLoading && !error && (
        <>
          {songs.length === 0 ? (
            <p>Список песен пуст</p>
          ) : (
            <>
              <div className="songs-grid">
                {songs.map((song) => (
                  <article className="song-card" key={song.id}>
                    <div className="song-card__icon">♫</div>

                    <div className="song-card__content">
                      {editingId === song.id ? (
                        <>
                          <input
                            type="text"
                            value={editName}
                            onChange={(event) => setEditName(event.target.value)}
                            placeholder="Название"
                          />

                          <input
                            type="text"
                            value={editArtist}
                            onChange={(event) =>
                              setEditArtist(event.target.value)
                            }
                            placeholder="Исполнитель"
                          />

                          <input
                            type="text"
                            value={editGenre}
                            onChange={(event) => setEditGenre(event.target.value)}
                            placeholder="Жанр"
                          />

                          <input
                            type="number"
                            min="0"
                            value={editDuration}
                            onChange={(event) =>
                              setEditDuration(event.target.value)
                            }
                            placeholder="Длительность"
                          />

                          <div className="song-card__actions">
                            <button
                              type="button"
                              onClick={() => updateSong(song.id)}
                            >
                              Сохранить
                            </button>

                            <button type="button" onClick={cancelEditing}>
                              Отмена
                            </button>
                          </div>
                        </>
                      ) : (
                        <>
                          <h2>{song.name}</h2>

                          <p className="song-card__artist">{song.artist}</p>

                          <div className="song-card__info">
                            <span>{song.genre}</span>

                            <span>
                              {song.duration !== null
                                ? formatDuration(song.duration)
                                : "Не указана"}
                            </span>
                          </div>

                          <div className="song-card__actions">
                            <button
                              type="button"
                              onClick={() => startEditing(song)}
                            >
                              Изменить
                            </button>

                            <button
                              type="button"
                              onClick={() => deleteSong(song.id)}
                            >
                              Удалить
                            </button>
                          </div>
                        </>
                      )}
                    </div>
                  </article>
                ))}
              </div>

              <div className="pagination">
                <button
                  type="button"
                  onClick={goToPreviousPage}
                  disabled={page === 1}
                >
                  Назад
                </button>

                <span>Страница {page}</span>

                <button type="button" onClick={goToNextPage}>
                  Вперёд
                </button>
              </div>

              {pageMessage && <p>{pageMessage}</p>}
            </>
          )}
        </>
      )}
    </main>
  )
}

export default Main