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
import { Link } from "react-router-dom";

const formSchema = z.object({
    email: z.string().email({
        message: "Please enter a valid email.",
    }),
    password: z.string().min(8, {
        message: "Password must be at least 10 characters.",
    }),
});

export const LoginPage = () => {
    const form = useForm<z.infer<typeof formSchema>>({
        defaultValues: {
            email: "",
            password: "",
        },
        resolver: zodResolver(formSchema),
    });

    function onSubmit(/*values: z.infer<typeof formSchema>*/) {
        //     signIn("credentials", {
        //       email: values.email as string,
        //       password: values.password as string
        //     });
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
