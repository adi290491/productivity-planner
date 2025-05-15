import axios from "axios";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import CarouselHeader from "../components/CarouselHeader";

const Signup = () => {
    const [form, setForm] = useState({name: '', email: '', password:''});
    const navigate = useNavigate();
    const handleChange = (field: string, value: string) => {
        setForm(prev => ({...prev, [field]: value}));
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        console.log("Signup form:", form)
        try {
            const response = await axios.post("http://localhost:8000/users/signup", form)
            alert(`Signup successful! Welcome ${response.data.user.name}`)
            navigate("/login")
        } catch (error: any) {
            alert("Signup failed: " + (error.response?.data?.error || error.message))
        }
    };

    return (
        <div className="bg-background min-h-screen">
            <CarouselHeader />
            <div className="max-w-lg mx-auto mt-12 p-8 bg-card border border-border shadow rounded">
                <h2 className="text-2xl font-semibold mb-6 text-text">Create an account</h2>
                <form className="space-y-5" onSubmit={handleSubmit}>
                <div>
                    <label className="block mb-1 text-sm text-accent">Name</label>
                    <input
                    type="text"
                    value={form.name}
                    onChange={e => handleChange('name', e.target.value)}
                    className="w-full px-4 py-2 border border-border rounded bg-white text-text"
                    />
                </div>
                <div>
                    <label className="block mb-1 text-sm text-accent">Email</label>
                    <input
                    type="email"
                    value={form.email}
                    onChange={e => handleChange('email', e.target.value)}
                    className="w-full px-4 py-2 border border-border rounded bg-white text-text"
                    />
                </div>
                <div>
                    <label className="block mb-1 text-sm text-accent">Password</label>
                    <input
                    type="password"
                    value={form.password}
                    onChange={e => handleChange('password', e.target.value)}
                    className="w-full px-4 py-2 border border-border rounded bg-white text-text"
                    />
                </div>
                <button
                    type="submit"
                    className="w-full bg-accent text-white py-2 rounded hover:brightness-110 transition"
                >
                    Sign Up
                </button>
                <p className="block mb-1 text-sm text-accent">
                    Already have an account? 
                    <span className="ml-1 underline cursor-pointer"
                    onClick={() => {
                        navigate("/login")
                    }}
                    >Sign In</span>
                </p>
                </form>
            </div>
        </div>
    );
};

export default Signup;