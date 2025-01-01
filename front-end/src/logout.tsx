import { Box, Container } from "@mui/material";
import { NavLink } from "react-router";

export default function Logout() {
  return (
    <div className="logut">
      <Container>
        <Box>
          <NavLink to={"/signin"}> Signin </NavLink>
          <p> bye bye</p>
        </Box>
      </Container>
    </div>
  );
}
