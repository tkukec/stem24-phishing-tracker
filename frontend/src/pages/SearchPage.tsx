import { EventsService } from "@/api/events";
import PhishingEventCard from "@/components/PhishingEventCard";
import { Button } from "@/components/ui/button";
import { FormField, FormItem, FormControl, Form } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import {
    Pagination,
    PaginationContent,
    PaginationItem,
    PaginationNext,
    PaginationPrevious,
} from "@/components/ui/pagination";
import { IPhishingEvent, PhishingEventSearchData } from "@/interfaces/PhishingEventIntefaces";
import Navbar from "@/layouts/Navbar";
import { zodResolver } from "@hookform/resolvers/zod";
import {
    Accordion,
    AccordionItem,
    AccordionTrigger,
    AccordionContent,
} from "@radix-ui/react-accordion";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

const events: Array<IPhishingEvent> = [
    {
        id: 1,
        name: "Event 1",
        createdAt: new Date(),
        brand: "Brand 1",
        description: "Event 1 description",
        maliciousUrl: "test",
        domainRegistrationDate: new Date(),
        keyword: [],
        status: "todo",
        comments: [],
        dnsRecords: [],
    },
    {
        id: 2,
        name: "Event 2",
        createdAt: new Date(),
        brand: "Brand 2",
        description:
            "Event 2 description Lorem ipsum dolor, sit amet consectetur adipisicing elit. Minus, ea odio atque voluptatibus laudantium itaque dolorum modi quidem nesciunt officia.",
        maliciousUrl: "",
        domainRegistrationDate: new Date(),
        keyword: ["keyword1", "short", "very long and cool", "keyword4"],
        status: "todo",
        comments: [],
        dnsRecords: [],
    },
];

const searchSchema = z.object({
    name: z.string().optional(),
    startDate: z.string().optional(),
    endDate: z.string().optional(),
    brand: z.string().optional(),
    domainName: z.string().optional(),
    keywords: z.string().optional(),
});
const SearchPage = () => {
    const [currentSearchData, setCurrentSearchData] = useState<PhishingEventSearchData>({
        name: "",
        startDate: new Date(),
        endDate: new Date(),
        brand: "",
        domainName: "",
        keywords: [],
    });
    const query = EventsService.useGetEvents(currentSearchData);
    const form = useForm<z.infer<typeof searchSchema>>({
        defaultValues: {
            name: "",
            startDate: "",
            endDate: "",
            brand: "",
            domainName: "",
            keywords: "",
        },
        resolver: zodResolver(searchSchema),
    });

    const onSubmit = (data: z.infer<typeof searchSchema>) => {
        setCurrentSearchData({
            ...data,
            startDate: data.startDate ? new Date(data.startDate) : new Date(),
            endDate: data.endDate ? new Date(data.endDate) : new Date(),
            keywords: data.keywords ? data.keywords?.split(/[, ]/) : [],
        });
    };

    return (
        <>
            <Navbar />
            <div className="flex justify-center py-10 px-4">
                <div className="flex flex-col gap-6">
                    <h2>Search page</h2>
                    <Form {...form}>
                        <form onSubmit={form.handleSubmit(onSubmit)}>
                            <Accordion type="single" collapsible className="w-full">
                                <AccordionItem value="item-1">
                                    <div className="flex gap-5">
                                        <FormField
                                            control={form.control}
                                            name="name"
                                            render={({ field }) => (
                                                <FormItem className="w-full">
                                                    <FormControl>
                                                        <Input
                                                            placeholder="Search by name"
                                                            {...field}
                                                        />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                        <AccordionTrigger className="text-sm flex gap-2" asChild>
                                            <Button variant="outline">Advanced</Button>
                                        </AccordionTrigger>
                                        <Button
                                            type="submit"
                                            className="bg-theme hover:opacity-60 hover:bg-theme"
                                        >
                                            Search
                                        </Button>
                                    </div>
                                    <AccordionContent className="flex flex-col gap-2 mt-4">
                                        <FormField
                                            control={form.control}
                                            name="brand"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <Input placeholder="Brand" {...field} />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="domainName"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <Input
                                                            placeholder="Domain Name"
                                                            {...field}
                                                        />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="keywords"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <Input placeholder="Keywords" {...field} />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="endDate"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <Input
                                                            type="date"
                                                            placeholder="End Date"
                                                            {...field}
                                                        />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                        <FormField
                                            control={form.control}
                                            name="startDate"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormControl>
                                                        <Input
                                                            type="date"
                                                            placeholder="Start Date"
                                                            {...field}
                                                        />
                                                    </FormControl>
                                                </FormItem>
                                            )}
                                        />
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>
                        </form>
                    </Form>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                        {events.map((event) => (
                            <PhishingEventCard key={event.id} event={event} />
                        ))}
                        {events.map((event) => (
                            <PhishingEventCard key={event.id} event={event} />
                        ))}
                        {events.map((event) => (
                            <PhishingEventCard key={event.id} event={event} />
                        ))}
                    </div>
                    <Pagination>
                        <PaginationContent>
                            {query.hasPreviousPage && (
                                <PaginationItem>
                                    <PaginationPrevious href="#" />
                                </PaginationItem>
                            )}
                            {query.hasNextPage && (
                                <PaginationItem>
                                    <PaginationNext href="#" />
                                </PaginationItem>
                            )}
                        </PaginationContent>
                    </Pagination>
                </div>
            </div>
        </>
    );
};

export default SearchPage;
