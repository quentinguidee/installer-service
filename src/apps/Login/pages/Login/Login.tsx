import "./Login.sass";
import {
    Button,
    Horizontal,
    Logo,
    MaterialIcon,
    TextField,
    Title,
    Vertical,
} from "@vertex-center/components";
import Spacer from "../../../../components/Spacer/Spacer";
import { APIError } from "../../../../components/Error/APIError";
import { useState } from "react";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useLogin } from "../../hooks/useLogin";

export default function Login() {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");

    const { login, isLoggingIn, errorLogin } = useLogin({
        onSuccess: () => {},
    });

    const onRegister = () => login({ username, password });
    const onUsernameChange = (e: any) => setUsername(e.target.value);
    const onPasswordChange = (e: any) => setPassword(e.target.value);

    return (
        <div className="login">
            <div className="login-container">
                <ProgressOverlay show={isLoggingIn} />
                <Horizontal gap={12}>
                    <Logo />
                    <Title variant="h1">Login</Title>
                </Horizontal>
                <Vertical gap={20}>
                    <TextField
                        id="username"
                        label="Username"
                        onChange={onUsernameChange}
                        required
                    />
                    <TextField
                        id="password"
                        label="Password"
                        onChange={onPasswordChange}
                        type="password"
                        required
                    />
                    <APIError error={errorLogin} />
                    <Horizontal>
                        <Spacer />
                        <Button
                            variant="colored"
                            rightIcon={<MaterialIcon icon="login" />}
                            onClick={onRegister}
                        >
                            Login
                        </Button>
                    </Horizontal>
                </Vertical>
            </div>
        </div>
    );
}
