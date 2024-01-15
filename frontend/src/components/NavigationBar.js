import useAuth from "../hooks/useAuth";
import { Nav, Navbar, Container } from "react-bootstrap";
import NotificationNavbarItem from "./NotificationNavbarItem";
import { LinkContainer } from "react-router-bootstrap";
import SearchUtility from "../components/SearchUtility";
import { BoxArrowRight } from "react-bootstrap-icons";
import NavbarChat from "../components/NavbarChat";
import NavbarGroupSidebar from "./NavbarGroupSidebar";
import SearchSmallUtility from "../components/SearchSmallUtility";

const NavigationBar = () => {
  const { auth } = useAuth();

  return (
    <>
      <Navbar
        collapseOnSelect
        expand="md"
        className="bg-body-tertiary"
        fixed="top"
      >
        <Container>
          <LinkContainer to="/">
            <Navbar.Brand>Social Hub</Navbar.Brand>
          </LinkContainer>
          {auth && (
            <>
              <NotificationNavbarItem />

              <div className="d-md-none">
                <NavbarChat />
              </div>
            </>
          )}
          <Navbar.Toggle aria-controls="responsive-navbar-nav" />

          <Navbar.Collapse id="responsive-navbar-nav">
            <Nav className="d-flex justify-content-between ms-auto">
              {auth && (
                <>
                  <LinkContainer to="/profile">
                    <Nav.Link>Profile</Nav.Link>
                  </LinkContainer>

                  <div className="d-none d-md-block">
                    <SearchUtility />
                  </div>
                  <div className="d-md-none">
                    <SearchSmallUtility />
                  </div>

                  <div className="d-md-none">
                    <NavbarGroupSidebar />
                  </div>
                </>
              )}

              {!auth ? (
                <>
                  <LinkContainer to="/login">
                    <Nav.Link>Sign In</Nav.Link>
                  </LinkContainer>
                  <LinkContainer to="/signup">
                    <Nav.Link>Sign Up</Nav.Link>
                  </LinkContainer>
                </>
              ) : (
                <LinkContainer to="/logout">
                  <Nav.Link>
                    <BoxArrowRight
                      className="d-flex align-self-center"
                      size={20}
                    />
                  </Nav.Link>
                </LinkContainer>
              )}
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
    </>
  );
};

export default NavigationBar;
