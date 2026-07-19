import React, { useEffect, useState } from "react";
import { type Project } from "../types/Project";
import { requireAuth } from "../utils/Auth";
import { useAuth } from "./AuthContext";
import type { AllProjectsDTO } from "../types/ProjectDTO";
import { useNavigate } from "react-router";

const BASE_URL = import.meta.env.VITE_API_URL;

function Dashboard() {
    const [projects, setProjects] = useState<Project[]>([])
    const [newProjectFormData, setNewProjectFormData] = useState({name: "", description: ""})
    const authContext = useAuth()
    const navigate = useNavigate();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setNewProjectFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const loadProjects = () => {
        const status = requireAuth(authContext)
        fetch(BASE_URL + "/api/v0/projects", {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + status.jwt
            }
        }).then(async (resp) => {
            if(resp.ok) {
                const data = await resp.json() as AllProjectsDTO
                console.log(data.data)
                setProjects(data.data.projects)
            }
        })
    }

    useEffect(() => {
        loadProjects();
    }, [])

    const addProject = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            
            const status = requireAuth(authContext)
            const resp = await fetch(BASE_URL + "/api/v0/projects", {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + status.jwt
                },
                body: JSON.stringify(newProjectFormData)
            })

            if(resp.ok) {
                // Reload all projects
                setNewProjectFormData({name: "", description: ""})
                loadProjects()
                console.log(await resp.json())
            } else {
                alert("Failed to create project!")
            }
        } catch (e) {
            console.log(e)
            alert("Failed to create project")
        }
    }

    return (
        <>
            <h1>This is a dashboard!</h1>
            <form id="create-project" onSubmit={addProject}>
                <input id="proj-name-input"
                    value={newProjectFormData.name}
                    onChange={handleChange}
                    type="text"
                    name="name"
                    placeholder="Enter project name" />
                <input id="proj-description-input"
                    value={newProjectFormData.description}
                    onChange={handleChange}
                    type="text"
                    name="description"
                    placeholder="Enter project description" />
                <button type="submit">Submit</button>
            </form>
            <hr></hr>
            {projects.map((val) => (
                <div key={val.id} onClick={() => navigate("/projects", {state: {project: val}})}>
                    <p >Project name:{val.name} Description:{val.description}</p>
                    <hr></hr>
                </div>
            ))}
        </>
    )
}

export default Dashboard;