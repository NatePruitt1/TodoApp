export interface Project {
    id: string;
    user_id: string;
    name: string;
    description: string;
    categories: Category[]
}

export interface Category {
    id: string;
    project_id: string;
    name: string;
    index: number;
    cards: Card[];
}

export interface Card {
    id: string;
    category_id: string;
    title: string;
    content: string;
    finished: boolean;
}