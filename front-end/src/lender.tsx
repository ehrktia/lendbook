import { Container, Box } from "@mui/material";
import { styled } from '@mui/material/styles';
import Paper from '@mui/material/Paper';
import { NavLink } from "react-router";
import Grid from '@mui/material/Grid2';
import Card from '@mui/material/Card';
import CardContent from '@mui/material/CardContent';
import CardMedia from '@mui/material/CardMedia';
import Typography from '@mui/material/Typography';
import CardActions from "@mui/material/CardActions";
import Button from "@mui/material/Button";

const Item = styled(Paper)(({ theme }) => ({
  backgroundColor: '#fff',
  ...theme.typography.body2,
  padding: theme.spacing(1),
  textAlign: 'center',
  color: theme.palette.text.secondary,
  ...theme.applyStyles('dark', {
    backgroundColor: '#1A2027',
  }),
}));

// eslint-disble-next-line
const backgroundURL = "http://localhost:8080/query";

export default function Lender() {
  return (
    <div className="lender">
      <Container>
        <Box sx={{ flexGrow: 1 }}>
          <Grid container spacing={4}>
            <Grid size={8}></Grid>
            <Grid size={4}>
              <Item>
                <ul>
                  <p> profile: first.last@email.com</p>
                  <NavLink to={"/logout"}> logout</NavLink>
                </ul>
              </Item>
            </Grid>
          </Grid>
        </Box>
        <Grid paddingY={2}></Grid>
        <Box sx={{ flexGrow: 1 }}>
          <Grid container spacing={4}>
            <Grid size={3}>
              <Item>
                <Card sx={{ maxWidth: 345 }}>
                  <CardMedia
                    sx={{ height: 150 }}
                    image="/book_1_cover.jpg"
                    title="Neelam-stupid book"
                  />
                  <CardContent>
                    <Typography gutterBottom variant="h6" component="div">
                      Neelam-stupid book-1
                    </Typography>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      This book takes an existing physics theory, fits in to 
                      real life with scenarios
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Box sx={{ flexGrow: 1 }}>
                      <Button size="small">Edit</Button>
                      <Button size="small">Remove</Button>
                    </Box>
                  </CardActions>
                </Card>
              </Item>
            </Grid>
            <Grid size={3}>
              <Item>
                <Card sx={{ maxWidth: 345 }}>
                  <CardMedia
                    sx={{ height: 150 }}
                    image="/book_2_cover.jpg"
                    title="Neelam stupid book 2"
                  />
                  <CardContent>
                    <Typography gutterBottom variant="h6" component="div">
                      Neelam stupid book version-2
                    </Typography>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      This book takes an existing physics theory and update it
                      with real life scenarios
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Box sx={{ flexGrow: 1 }}>
                      <Button size="small">Edit</Button>
                      <Button size="small">Remove</Button>
                    </Box>
                  </CardActions>
                </Card>
              </Item>
            </Grid>
            <Grid size={3}>
              <Item>
                <Card sx={{ maxWidth: 345 }}>
                  <CardMedia
                    sx={{ height: 150 }}
                    image="/book_3_cover.jpg"
                    title="neelam stupid book-3.0"
                  />
                  <CardContent>
                    <Typography gutterBottom variant="h6" component="div">
                      Neelam stupid book 3.0
                    </Typography>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      This is 3rd version of neelam stupid book, enhanced 
                      examples of the life matching physics theory
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Box sx={{ flexGrow: 1 }}>
                      <Button size="small">Edit</Button>
                      <Button size="small">Remove</Button>
                    </Box>
                  </CardActions>
                </Card>
              </Item>
            </Grid>
            <Grid size={3}>
              <Item>
                <Card sx={{ maxWidth: 345 }}>
                  <CardMedia
                    sx={{ height: 150 }}
                    image="/lizard.jpg"
                    title="neelam book on astrology"
                  />
                  <CardContent>
                    <Typography gutterBottom variant="h6" component="div">
                      Neelam Astrology Book-4
                    </Typography>
                    <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                      This book explains why saturn rules your life and not 
                      free will
                    </Typography>
                  </CardContent>
                  <CardActions>
                    <Box sx={{ flexGrow: 1 }}>
                      <Button size="small">Edit</Button>
                      <Button size="small">Remove</Button>
                    </Box>
                  </CardActions>
                </Card>
              </Item>
            </Grid>
          </Grid>
        </Box>
      </Container>
    </div >
  );
};

