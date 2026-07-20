import { useCardsApi } from "../api/cards";
import type { Card } from "../types/Project";

export function CardElement({card, reloadProject}: {card: Card, reloadProject: {(): void}}) {
    const cardsApi = useCardsApi();

    console.log(card.finished)

    const handleDelete = async () => {
        try {
            await cardsApi.remove(card.id)
            reloadProject()
        } catch(error) {
            console.error(error)
            alert("Failed to delete card.")
        }
    }

    return (
        <div className="card">
            <h4 contentEditable={true}>{card.title}</h4>
            <p contentEditable={true}>{card.content}</p>
            <button onClick={handleDelete}>Delete</button>
        </div>
    )
}