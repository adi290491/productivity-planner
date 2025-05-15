import { useLocation, useNavigate } from "react-router-dom";


const CarouselHeader = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const token = localStorage.getItem("token");

    const handleLogout = () => {
        localStorage.removeItem("token");
        navigate("/login")
    };

    const isAuthPage = location.pathname === '/login' || location.pathname === '/signup';

    return (
        <header className="w-full bg-primary text-text border-b border-border shadow-sm">
            <div className="max-w-screen-xl mx-auto px-8 py-6 flex flex-col md:flex-row justify-between items-start md:items-center">
                <div>
                    <h1 className="text-4xl font-bold">Productivity Tracker</h1>
                    <p className="text-accent text-sm italic mt-2 md:mt-0">with trend analysis</p>
                </div>
                {!isAuthPage && token && (
                    <button
                    onClick={handleLogout}
                    className="mt-4 md:mt-0 bg-accent text-white px-4 py-2 rounded hover:brightness-110 transition"
                    >
                        Logout
                    </button>
                )}
            </div>
        </header>
    );
}

export default CarouselHeader;