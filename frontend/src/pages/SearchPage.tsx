import PhishingEventCard from "@/components/PhishingEventCard";
import {Input} from "@/components/ui/input";
import {
    Pagination,
    PaginationContent,
    PaginationEllipsis,
    PaginationItem,
    PaginationLink,
    PaginationNext,
    PaginationPrevious,
} from "@/components/ui/pagination";
import {IPhishingEvent} from "@/interfaces/PhishingEventIntefaces";
import Navbar from "@/layouts/Navbar.tsx";

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

const SearchPage = () => {
    return (
        <>
            <Navbar/>
            <div className="flex flex-col gap-6">
                <h2>Search page</h2>
                <Input placeholder="Search"/>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                    {events.map((event) => (
                        <PhishingEventCard key={event.id} event={event}/>
                    ))}
                    {events.map((event) => (
                        <PhishingEventCard key={event.id} event={event}/>
                    ))}
                    {events.map((event) => (
                        <PhishingEventCard key={event.id} event={event}/>
                    ))}
                </div>
                <Pagination>
                    <PaginationContent>
                        <PaginationItem>
                            <PaginationPrevious href="#"/>
                        </PaginationItem>
                        <PaginationItem>
                            <PaginationLink href="#">1</PaginationLink>
                        </PaginationItem>
                        <PaginationItem>
                            <PaginationLink href="#">2</PaginationLink>
                        </PaginationItem>
                        <PaginationItem>
                            <PaginationEllipsis/>
                        </PaginationItem>
                        <PaginationItem>
                            <PaginationNext href="#"/>
                        </PaginationItem>
                    </PaginationContent>
                </Pagination>
            </div>
        </>
    );
};

export default SearchPage;
