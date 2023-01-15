class Home extends React.Component {

    constructor(props) {
        super(props);
        this.requestBoard = this.requestBoard.bind(this);
    }

    requestBoard() {
        console.log("Requesting Board")
        const getBoard = () => {
            fetch("http://localhost:3399/api/board").then((res) => res.json()).then((json) => { console.log(json); })
        }
        getBoard()
}

componentDidMount() {
    this.requestBoard();
}

render() {
    return (
        <div className="container">
            <div className="jumbotron">
                <h1>Godoku</h1>
                {/* <p>
                        {this.state.board.map(function(cell, i){
                            return <Cell key={i} cell={cell} />;
                        })}
                    </p> */}
            </div>
        </div>
    )
}
}

class Cell extends React.Component {
    constructor(props) {
        super(props)

    }
}
class App extends React.Component {

    render() {
        return (<Home />);
    }

}

ReactDOM.render(<App />, document.getElementById('app'));