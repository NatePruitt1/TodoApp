import { useState } from "react";
import type { Card, Category } from "../types/Project";
import { useAuth } from "./AuthContext";
import { requireAuth } from "../utils/Auth";

const BASE_URL = import.meta.env.VITE_API_URL

export function CategoryCard({category, reloadProject}: {category: Category, reloadProject: {(): void}}) {
    const authContext = useAuth()
    const [formData, setFormData] = useState({title: "", content: ""})
    
    const addCard = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const status = requireAuth(authContext)
            const resp = await fetch(BASE_URL + "/api/v0/categories/" + category.id + "/cards", {
                method: "POST",
                body: JSON.stringify({title: formData.title, content: formData.content}),
                headers: {
                    "Authorization": "Bearer " + status.jwt
                }
            })

            if(resp.ok) {
                setFormData({title: "", content: ""})
                reloadProject()
            } else {
                alert("Could not create card.")
            }
        } catch(e) {
            alert("Could not create card. ERROR ")
            console.log(e)
        }
    }
    
    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const moveCategory = async () => {
        const index = category.index + 1;
        const status = requireAuth(authContext)
        const resp = await fetch(BASE_URL + "/api/v0/categories/" + category.id + "/position", {
            method: "PATCH",
            body: JSON.stringify({index: index}),
            headers: {
                "Authorization": "Bearer " + status.jwt
            }
        })

        if (resp.ok) {
            reloadProject();
        } else {
            console.log("Wrong")
        }
    }

    const moveCategoryBack = async () => {
        const index = category.index - 1;
        const status = requireAuth(authContext)
        const resp = await fetch(BASE_URL + "/api/v0/categories/" + category.id + "/position", {
            method: "PATCH",
            body: JSON.stringify({index: index}),
            headers: {
                "Authorization": "Bearer " + status.jwt
            }
        })

        if (resp.ok) {
            reloadProject();
        } else {
            console.log("Wrong")
        }
    }

    const deleteCategory = async () => {
        const status = requireAuth(authContext)
        const resp = await fetch(BASE_URL + "/api/v0/categories/" + category.id, {
            method: "DELETE",
            headers: {
                "Authorization": "Bearer " + status.jwt
            }
        })

        if(resp.ok) {
            reloadProject();
        }
    }

    return (
        <div key={category.id} className="category">
            <hr />
            <h3>{category.name}</h3>
            <button id="move-forward" onClick={moveCategory}>&gt;</button>
            <button id="move-backward" onClick={moveCategoryBack}>&lt;</button>
            <button id="delete" onClick={deleteCategory}>Delete</button>
            <form id="form-create-card" onSubmit={addCard}>
                <input id="card-name-input"
                    value={formData.title}
                    onChange={handleChange}
                    type="text"
                    name="title"
                    placeholder="Enter card title." />
                <input id="card-content-input"
                    value={formData.content}
                    onChange={handleChange}
                    type="text"
                    name="content"
                    placeholder="Enter card content." />
                <button type="submit">Add card</button>
            </form>
            {category.cards.map((v: Card) => (
                <div key={v.id}>
                    <hr/>
                    <p>Title: {v.title}</p>
                    <p>Content: {v.content}</p>
                    <hr/>
                </div>
            ))}
            <hr />
        </div>
    )
}