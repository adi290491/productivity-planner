import { useState } from "react";
import { useNavigate } from "react-router-dom";
import CarouselHeader from "../components/CarouselHeader";
import axios from "axios";

const Login = () => {
    const [form, setForm] = useState({email: '', password: ''});
    const navigate = useNavigate()

    const handleChange = (field: string, value: string) => {
        setForm(prev => ({...prev, [field]:value}) )
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        try {
            const response = await axios.post("http://localhost:8000/users/login", form)
            localStorage.setItem("token", response.data.token)
            console.log('Response:', response)
            // alert(`Signin successful! Welcome ${response.data.user.name}`)
            navigate("/")
        } catch (error: any) {
             alert("Signin failed: " + (error.response?.data?.error || error.message))
        }
    };

    return (
        <div className="bg-background min-h-screen">
            <CarouselHeader />
            <div className="max-w-lg mx-auto mt-12 p-8 bg-card border border-border shadow rounded">
                <h2 className="text-2xl font-semibold mb-6 text-text">User Login</h2>
                <form className="space-y-5" onSubmit={handleSubmit}>
                <div>
                    <label className="block mb-1 text-sm text-accent">Email</label>
                    <input
                    type="email"
                    required
                    value={form.email}
                    onChange={e => handleChange('email', e.target.value)}
                    className="w-full px-4 py-2 border border-border rounded bg-white text-text"
                    />
                </div>
                <div>
                    <label className="block mb-1 text-sm text-accent">Password</label>
                    <input
                    type="password"
                    required
                    value={form.password}
                    onChange={e => handleChange('password', e.target.value)}
                    className="w-full px-4 py-2 border border-border rounded bg-white text-text"
                    />
                </div>
                <button
                    type="submit"
                    className="w-full bg-accent text-white py-2 rounded hover:brightness-110 transition"
                >
                    Sign In
                </button>
                <p className="block mb-1 text-sm text-accent">
                    Don't have an account? 
                    <span className="ml-1 underline cursor-pointer"
                    onClick={() => {
                        navigate("/signup")
                    }}
                    >Create Account</span>
                </p>
                </form>
            </div>
        </div>
    );
};

export default Login;