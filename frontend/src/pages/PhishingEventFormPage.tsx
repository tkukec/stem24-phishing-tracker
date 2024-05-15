import PhishingEventForm from "@/components/Forms/PhishingEventForm.tsx";
import Navbar from "@/layouts/Navbar.tsx";


const PhishingEventFormPage = () => {
    return (
        <div>
            <Navbar/>
            <div className="flex justify-center items-center">
                <PhishingEventForm/>
            </div>
        </div>
    );
};

export default PhishingEventFormPage;
