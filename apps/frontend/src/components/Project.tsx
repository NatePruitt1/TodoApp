import { useLocation, useNavigate } from "react-router";
import type { Project } from "../types/Project";
import { useEffect, useState } from "react";
import { CategoryCard } from "./Category";
import { useProjectsApi } from "../api/projects";
import { useCategoriesApi } from "../api/categories";
import "./Project.css"

import addUrl from "../../public/material-symbols--add.svg"

interface ProjectState {
    project: Project
}

export function ProjectScreen() {
    const projectApi = useProjectsApi();
    const categoriesApi = useCategoriesApi();

    const location = useLocation();
    const navigate = useNavigate();
    const state: ProjectState = location.state || {}
    const [formData, setFormData]= useState({name: ""});
    
    const [project, setProject] = useState<Project>(state.project);

    const getProject = async () => {
        try {
            const proj = await projectApi.get(state.project.id)
            if(proj.categories && proj.categories.length > 0) {
                proj.categories.sort((a, b) => a.index - b.index)
            }
            setProject(proj)
        } catch (error) {
            console.error(error)
            alert("Could not load project from api.")
            navigate("/dashboard", {replace: true})
        }
    }

    useEffect(() => {
        getProject()
    }, [])

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };
    
    const addCategory = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            await categoriesApi.create(project.id, {name: formData.name})
            setFormData({name: ""})
            getProject()
        } catch (e) {
            console.error(e)
            alert("Failed to create project")
        }
    }

    return (
        
        <div className="project">
            <header className="header-bar">
                <h3>Todo - Project</h3>
            </header>
            <div className="project-header">
                <h1>{project?.name}</h1>
                <p>- {project?.description}</p>
            </div>
            <div className="categories">
                {project?.categories?.map((category) => (
                    <CategoryCard key={category.id} category={category} reloadProject={getProject} />
                ))}
                <div className="category">
                    <div className="create-category-content">
                        <p className="category-title">Add Category</p>
                        <img className="small-icon" src={addUrl} />
                    </div>
                </div>
            </div>
        </div>
        
    )
}