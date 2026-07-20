import React, { useEffect, useState } from "react";
import { type Project } from "../types/Project";
import { useNavigate } from "react-router";
import { useProjectsApi } from "../api/projects";

function Dashboard() {
    const [projects, setProjects] = useState<Project[]>([])
    const [newProjectFormData, setNewProjectFormData] = useState({name: "", description: ""})
    const navigate = useNavigate();
    const projectApi = useProjectsApi();

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setNewProjectFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };

    const loadProjects = async () => {
        try {
            const projects = await projectApi.list()
            setProjects(projects.projects)
        } catch(error) {
            console.error(error)
            alert("Failed to load projects.")
        }
    }

    useEffect(() => {
        loadProjects();
    }, [])

    const addProject = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            await projectApi.create({name: newProjectFormData.name, description: newProjectFormData.description})
            setNewProjectFormData({name: "", description: ""})
            loadProjects()
        } catch (e) {
            console.error(e)
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