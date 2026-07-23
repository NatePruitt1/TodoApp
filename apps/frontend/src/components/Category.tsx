import type { Card, Category } from "../types/Project";
import { useCategoriesApi } from "../api/categories";
import { CardElement } from "./Card";
import "./Category.css"

import addUrl from "../../public/material-symbols--add.svg"
import trashUrl from "../../public/material-symbols--delete-outline.svg"

export function CategoryCard({category, reloadProject}: {category: Category, reloadProject: {(): void}}) {
    const categoriesApi = useCategoriesApi();

    const deleteCategory = async () => {
        try{
            await categoriesApi.remove(category.id)
            reloadProject()
        } catch (error) {
            console.error(error)
            alert("Could not delete category.")
        }
    }

    const startDrag = (e: React.DragEvent<HTMLDivElement>) => {
        e.dataTransfer.setData("application/json", JSON.stringify(category))
    }

    const dragOver = (e: React.DragEvent<HTMLDivElement>) => {
        e.preventDefault()
    }

    const dropOver = (e: React.DragEvent<HTMLDivElement>) => {
        const data = JSON.parse(e.dataTransfer.getData('application/json')) as Category;
        console.log(`Move category index: ${data.index} to ${category.index}`);
        (async () => {
            try {
                await categoriesApi.reposition(data.id, category.index)
                reloadProject()
            } catch(error) {
                console.error(error)
            }
        })();
    }

    return (
        <div key={category.id} className="category" draggable={true} onDragStart={startDrag} onDragOver={dragOver} onDrop={dropOver} >
            <div className="category-header">
                <h3 className="category-title">{category.name}</h3>
                <button id="delete" className="small-icon-button" onClick={deleteCategory}><img src={trashUrl} /></button>
            </div>

            {category.cards.map((v: Card) => (
                <CardElement key={v.id} card={v} reloadProject={reloadProject} />
            ))}

            <div className="card">
                <div className="create-card-content">
                    <p>Add a card </p>
                    <img src={addUrl} className="medium-icon" />
                </div>
            </div>
        </div>
    )
}