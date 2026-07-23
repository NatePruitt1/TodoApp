import { useLocation, useNavigate } from "react-router";
import type { Project } from "../types/Project";
import { useEffect, useState } from "react";
import { CategoryCard } from "./Category";
import { useProjectsApi } from "../api/projects";
import { useCategoriesApi } from "../api/categories";
import "./Project.css"

import addUrl from "../../public/material-symbols--add.svg"
import { Header } from "./Header";

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
            <Header title="Todo - Project" />
            <div className="project-header">
                <h4>{project?.name}</h4>
                <p>- {project?.description}</p>
            </div>
            <div className="categories">
                {project?.categories?.map((category) => (
                    <CategoryCard key={category.id} category={category} reloadProject={getProject} />
                ))}
                <div className="category">
                    <div className="create-category-content">
                        <p className="category-title">Add Category</p>
                    </div>
                    <form id="create-category-form" onSubmit={addCategory}>
                        <input id="dashboard-name-input"
                            value={formData.name}
                            onChange={handleChange}
                            type="text"
                            name="name" 
                            placeholder="Enter Your Category's Name." 
                            autoComplete="off" />
                        <button type="submit" className="tiny-icon-button"><img src={addUrl} /></button>
                    </form>
                </div>
            </div>
        </div>
        
    )
}