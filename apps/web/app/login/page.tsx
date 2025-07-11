import { LoginForm } from "@/components/login-form"

export default function LoginPage() {
    return (
        <main className="flex flex-col w-full min-h-svh items-center justify-center p-4">
            <section className="w-full max-w-sm">
                <LoginForm />
            </section>
        </main>
    )
}