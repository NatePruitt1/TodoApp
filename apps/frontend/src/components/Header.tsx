import { useNavigate } from "react-router"

import homeUrl from "../../public/icons8-home.svg"
import personUrl from "../../public/icons8-person-64.png"
import addUrl from "../../public/icons8-plus-50.png"
import { useState } from "react";
import { useProjectsApi } from "../api/projects";

export function Header({title}: {title: string}) {
    const navigate = useNavigate();
    const [newProjectFormData, setNewProjectFormData] = useState({name: "", description: ""})
    const projectApi = useProjectsApi();

    const addProject = async (event: React.SubmitEvent<HTMLFormElement>) => {
        event.preventDefault();
        try {
            const project = await projectApi.create({name: newProjectFormData.name, description: newProjectFormData.description})
            setNewProjectFormData({name: "", description: ""})
            navigate("/projects", {state: {project: project}})
        } catch (e) {
            console.error(e)
            alert("Failed to create project")
        }
    }
    
    const toggleCreateDialog = () => {
        const menu = document.getElementById('dashboard-create-dialog') as HTMLDialogElement
        menu.showModal()
    }

    const closeCreateDialog = () => {
        const menu = document.getElementById('dashboard-create-dialog') as HTMLDialogElement
        menu.close()
        setNewProjectFormData({name: "", description: ""})
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement> | React.ChangeEvent<HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setNewProjectFormData((prev) => ({
            ...prev,
            [name]: value,
        }));
    };
    
    return (
    <>
        <header className="header-bar">
            <div className="header-home-button-grouping">
                <button onClick={() => navigate("/")} className="small-icon-button"><img className="medium-icon" src={homeUrl} /></button>
                <h3>{title}</h3>
            </div>
            <div className="header-button-grouping">
                <button onClick={toggleCreateDialog} id="create-button" className="small-icon-button"><img className="medium-icon" src={addUrl} /></button>
                <button onClick={()=>navigate("/Profile")} className="small-icon-button"><img className="medium-icon" src={personUrl} /></button>
            </div>
        </header>

        <dialog id="dashboard-create-dialog">
            <div className="dashboard-create-header">
                <h3>New Project</h3>
                <button onClick={closeCreateDialog}>&times;</button>
            </div>
            <p>Create a new Project.</p>
            <form id="dashboard-create-form" onSubmit={addProject}>
                <label>Name</label>
                <input id="dashboard-name-input"
                    value={newProjectFormData.name}
                    onChange={handleChange}
                    type="text"
                    name="name" 
                    placeholder="Enter Your Projects Name."
                    className="dashboard-create-input" 
                    autoComplete="off" />
                <label>Description</label>
                <textarea id="dashboard-description-input"
                    value={newProjectFormData.description}
                    onChange={handleChange}
                    name="description"
                    placeholder="Describe Your Project."
                    className="dashboard-create-input dashboard-create-textarea"
                    autoComplete="off" />
                <button type="submit">Create</button>
            </form>
        </dialog>
        </>
    )
}