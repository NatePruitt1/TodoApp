import { useLocation } from "react-router";
import type { Project } from "../types/Project";
import { useEffect, useState } from "react";
import { requireAuth } from "../utils/Auth";
import { useAuth } from "./AuthContext";
import { CategoryCard } from "./Category";

const BASE_URL = import.meta.env.VITE_API_URL

interface ProjectState {
    project: Project
}

export function Project() {
    const authContext = useAuth();
    const location = useLocation();
    const state: ProjectState = location.state || {}
    const [formData, setFormData]= useState({name: ""});
    
    const [project, setProject] = useState<Project>();

    const getProject = async () => {
        const status = requireAuth(authContext)
        const resp = await fetch(BASE_URL + "/api/v0/projects/" + state.project.id, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + status.jwt
            }
        })

        if(resp.ok) {
            const respData = await resp.json()
            const proj = respData.data as Project
            if(proj.categories && proj.categories.length > 0) {
                proj.categories.sort((a, b) => a.index - b.index)
            }
            setProject(proj)
        } else {
            alert("Could not load project!")
            console.log(await resp.json())
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
            
            const status = requireAuth(authContext)
            const path = '/api/v0/projects/' + state.project.id + '/categories'
            const resp = await fetch(BASE_URL + path, {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + status.jwt
                },
                body: JSON.stringify(formData)
            })

            if(resp.ok) {
                // Reload all projects
                setFormData({name: ""})
                console.log(await resp.json())
                getProject()
            } else {
                alert("Failed to create project!")
            }
        } catch (e) {
            console.log(e)
            alert("Failed to create project")
        }
    }

    return (
    
        <div className="project">
            <form id="form-create-category" onSubmit={addCategory}>
                <input id="cat-name-input"
                    value={formData.name}
                    onChange={handleChange}
                    type="text"
                    name="name"
                    placeholder="Enter category name." />
                <button type="submit">Add category</button>
            </form>
            <h1>{project?.name}</h1>
            <p>{project?.description}</p>
            <div className="categories">
                {project?.categories.map((category) => (
                    <CategoryCard key={category.id} category={category} reloadProject={getProject} />
                ))}
            </div>
        </div>
        
    )
}