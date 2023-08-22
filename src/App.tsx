import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import Infrastructure from "./pages/Infrastructure/Infrastructure";
import BayDetails from "./pages/BayDetails/BayDetails";
import BayDetailsLogs from "./pages/BayDetailsLogs/BayDetailsLogs";
import BayDetailsEnv from "./pages/BayDetailsEnv/BayDetailsEnv";
import BayDetailsDependencies from "./pages/BayDetailsDependencies";
import BayDetailsHome from "./pages/BayDetailsHome/BayDetailsHome";
import Settings from "./pages/Settings/Settings";
import SettingsTheme from "./pages/SettingsTheme/SettingsTheme";
import { useContext } from "react";
import { ThemeContext } from "./main";
import classNames from "classnames";
import SettingsAbout from "./pages/SettingsAbout/SettingsAbout";
import SettingsUpdates from "./pages/SettingsUpdates/SettingsUpdates";
import BayDetailsDocker from "./pages/BayDetailsDocker/BayDetailsDocker";
import BayDetailsSettings from "./pages/BayDetailsSettings/BayDetailsSettings";
import BayDetailsStatus from "./pages/BayDetailsStatus/BayDetailsStatus";
import Store from "./pages/Store/Store";
import Dock from "./components/Dock/Dock";
import BayDetailsUpdate from "./pages/BayDetailsUpdate/BayDetailsUpdate";

function App() {
    const { theme } = useContext(ThemeContext);

    return (
        <div className={classNames("app", theme)}>
            <HashRouter>
                <Routes>
                    <Route
                        path="/"
                        element={<Navigate to="/infrastructure" />}
                        index
                    />
                    <Route
                        path="/infrastructure"
                        element={<Infrastructure />}
                    />
                    <Route
                        path="/infrastructure/:uuid/"
                        element={<BayDetails />}
                    >
                        <Route
                            path="/infrastructure/:uuid/"
                            element={<BayDetailsHome />}
                        />
                        <Route
                            path="/infrastructure/:uuid/docker"
                            element={<BayDetailsDocker />}
                        />
                        <Route
                            path="/infrastructure/:uuid/logs"
                            element={<BayDetailsLogs />}
                        />
                        <Route
                            path="/infrastructure/:uuid/status"
                            element={<BayDetailsStatus />}
                        />
                        <Route
                            path="/infrastructure/:uuid/environment"
                            element={<BayDetailsEnv />}
                        />
                        <Route
                            path="/infrastructure/:uuid/dependencies"
                            element={<BayDetailsDependencies />}
                        />
                        <Route
                            path="/infrastructure/:uuid/update"
                            element={<BayDetailsUpdate />}
                        />
                        <Route
                            path="/infrastructure/:uuid/settings"
                            element={<BayDetailsSettings />}
                        />
                    </Route>
                    <Route path="/settings" element={<Settings />}>
                        <Route
                            path="/settings/theme"
                            element={<SettingsTheme />}
                        />
                        <Route
                            path="/settings/updates"
                            element={<SettingsUpdates />}
                        />
                        <Route
                            path="/settings/about"
                            element={<SettingsAbout />}
                        />
                    </Route>
                    <Route path="/marketplace" element={<Store />} />
                </Routes>
                <Dock />
            </HashRouter>
        </div>
    );
}

export default App;
