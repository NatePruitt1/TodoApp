import { useState } from "react";
import type { Card, Category } from "../types/Project";
import { useCategoriesApi } from "../api/categories";
import { useCardsApi } from "../api/cards";
import { CardElement } from "./Card";

export function CategoryCard({category, reloadProject}: {category: Category, reloadProject: {(): void}}) {
    const categoriesApi = useCategoriesApi();
    const cardsApi = useCardsApi();

    const [formData, setFormData] = useState({title: "", content: ""})
    
    const addCard = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            await cardsApi.create(category.id, {title: formData.title, content: formData.content})
            setFormData({title: "", content: ""})
            reloadProject()
        } catch(e) {
            alert("Could not create card.")
            console.error(e)
        }
    }
    
    const handleChange = (e: React.ChangeEvent<HTMLInputElement> | React.ChangeEvent<HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const moveCategory = async () => {
        try {
            await categoriesApi.reposition(category.id, category.index + 1)
            reloadProject();
        } catch(error) {
            console.error(error)
            alert("Could not move category.")
        }
    }

    const moveCategoryBack = async () => {
        try {
            await categoriesApi.reposition(category.id, category.index - 1)
            reloadProject();
        } catch(error) {
            console.error(error)
            alert("Could not move category.")
        }
    }

    const deleteCategory = async () => {
        try{
            await categoriesApi.remove(category.id)
            reloadProject()
        } catch (error) {
            console.error(error)
            alert("Could not delete category.")
        }
    }

    return (
        <div key={category.id} className="category" draggable={true}>
            <h3 className="category-title">{category.name}</h3>
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
                <textarea id="card-content-input"
                    value={formData.content}
                    onChange={handleChange}
                    name="content"
                    placeholder="Enter card content." />
                <button type="submit">Add card</button>
            </form>
            {category.cards.map((v: Card) => (
                <CardElement key={v.id} card={v} reloadProject={reloadProject} />
            ))}
        </div>
    )
}