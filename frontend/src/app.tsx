import './app.css'
import {
    StaticRouter,
    Routes,
    Route
} from 'react-router-dom';
import Post from "./components/post.tsx";
import Landing from "./components/landing.tsx";
import Login from "./components/login.tsx";
import Signup from "./components/signup.tsx";
import Profile from "./components/profile.tsx";
import Journal from "./components/journal.tsx";

function App(url: string) {
    return (
        <StaticRouter location={ url }>
            <Routes>
                <Route path="/" Component={Landing} />
                <Route path="/login" Component={Login} />
                <Route path="/signup" Component={Signup} />
                <Route path="/profile" Component={Profile} />
                <Route path="/journal" Component={Journal} />
                <Route path="/post" Component={Post} />
            </Routes>
        </StaticRouter>
    );
}

export default App;