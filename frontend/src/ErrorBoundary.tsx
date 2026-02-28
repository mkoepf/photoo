import { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error?: Error;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error("Uncaught error:", error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="error-fallback">
                    <h1>Something went wrong.</h1>
                    <p>The application crashed during render.</p>
                    <pre style={{ textAlign: 'left', background: '#333', padding: '1rem', overflow: 'auto' }}>
                        {this.state.error?.toString()}
                        {"\n\n"}
                        {this.state.error?.stack}
                    </pre>
                    <button onClick={() => window.location.reload()}>Reload App</button>
                </div>
            );
        }

        return this.props.children;
    }
}
