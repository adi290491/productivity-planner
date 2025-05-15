import { Navigate } from "react-router-dom";
import type { JSX } from "react/jsx-runtime";

const PrivateRoute = ({children} : {children: JSX.Element}) => {
    const token = localStorage.getItem('token');
    return token ? children : <Navigate to="/login" replace />;
};

export default PrivateRoute;