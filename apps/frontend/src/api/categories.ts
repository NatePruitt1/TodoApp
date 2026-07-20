import type { Category } from "../types/Project";
import { useApi } from "./useApi";

export function useCategoriesApi() {
    const api = useApi();
    return {
        create: (projectId: string, body: {name: string}) =>
            api.post<Category>(`/projects/${projectId}/categories`, body),
        update: (categoryId: string, body: {name: string}) =>
            api.patch<Category>(`/categories/${categoryId}`, body),
        remove: (categoryId: string) =>
            api.delete<void>(`/categories/${categoryId}`),
        reposition: (categoryId: string, index: number) =>
            api.patch<Category>(`/categories/${categoryId}/position`, {index})
    }
}