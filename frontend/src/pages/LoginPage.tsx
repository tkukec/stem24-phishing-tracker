import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";
import { Button } from "@/components/ui/button";
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Link, useNavigate } from "react-router-dom";
import AuthService from "@/api/auth";

const formSchema = z.object({
    email: z.string().email({
        message: "Please enter a valid email.",
    }),
    password: z.string().min(1, {
        message: "Password is required",
    }),
});

export const LoginPage = () => {
    const navigate = useNavigate();
    const form = useForm<z.infer<typeof formSchema>>({
        defaultValues: {
            email: "",
            password: "",
        },
        resolver: zodResolver(formSchema),
    });
    const { mutate: login } = AuthService.useLogin();
    function onSubmit(values: z.infer<typeof formSchema>) {
        login(values, {
            onSuccess: () => {
                navigate("/");
            },
        });
    }

    return (
        <div className="w-full h-screen bg-theme overflow-hidden flex justify-center items-center">
            <div className="rounded-xl border bg-card text-card-foreground shadow p-8">
                <Form {...form}>
                    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-2 w-300">
                        <h1 className="heading3">Log in</h1>
                        <FormField
                            control={form.control}
                            name="email"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>E-mail</FormLabel>
                                    <FormControl>
                                        <Input {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="password"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel>{"Password"}</FormLabel>
                                    <FormControl>
                                        <Input type="password" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <Button className="w-full" type="submit">
                            {"Submit"}
                        </Button>
                        <Link to="/register">
                            <Button className="w-full" variant={"link"}>
                                Don't have an account? Register here
                            </Button>
                        </Link>
                    </form>
                </Form>
            </div>
        </div>
    );
};
