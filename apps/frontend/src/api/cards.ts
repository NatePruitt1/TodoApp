import { useApi } from "./useApi";

export function useCardsApi() {
    const api = useApi();
    return {
        create: (categoryId: string, body: {title: string, content?: string}) =>
            api.post(`/categories/${categoryId}/cards`, body),
        remove: (cardId: string) => api.delete(`/cards/${cardId}`),
        retitle: (cardId: string, title: string) => 
            api.patch(`/cards/${cardId}`, {title}),
        move: (cardId: string, categoryId: string) => api.patch(`/cards/${cardId}/move`, {category_id: categoryId}),
        finish: (cardId: string, finished: boolean) => api.patch(`/cards/${cardId}/finish`, {finished})
    }
}