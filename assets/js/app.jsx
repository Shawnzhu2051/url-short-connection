class App extends React.Component {
    render() {
        return(
            <Home />
        );
    }
}

class Home extends React.Component {
    constructor() {
        super();
        this.long2shortURL = "http://localhost:8080/long2short"
        this.short2longURL = "http://localhost:8080/short2long"
        this.state = {
            longURLValue: '',
            shortURLValue: '',
            longURLResult: '',
            shortURLResult: ''
        };
        this.handleLongURLChange = this.handleLongURLChange.bind(this);
        this.handleShortURLChange = this.handleShortURLChange.bind(this);
        this.handleLongURLSubmit = this.handleLongURLSubmit.bind(this);
        this.handleShortURLSubmit = this.handleShortURLSubmit.bind(this);
    }

    handleLongURLChange(event) {
        this.setState({longURLValue: event.target.value});
    }

    handleShortURLChange(event) {
        this.setState({shortURLValue: event.target.value});
    }

    handleLongURLSubmit(event) {
        let formData = new FormData();
        formData.append("longUrl", this.state.longURLValue);
    
        fetch(this.long2shortURL , {
            method: 'POST',
            headers: {},
            body: formData,
        }).then((response) => {
            if (response.ok) {
                return response.json();
            }}).then((json) => {
                this.setState({
                    shortURLResult: JSON.stringify(json.shortUrl)
                })
            }).catch((error) => {
                console.error(error);
            });
        event.preventDefault();
    }

    handleShortURLSubmit(event) {
        let formData = new FormData();
        formData.append("shortUrl", this.state.shortURLValue);
    
        fetch(this.short2longURL , {
            method: 'POST',
            headers: {},
            body: formData,
        }).then((response) => {
            if (response.ok) {
                return response.json();
            }}).then((json) => {
                this.setState({
                    longURLResult: JSON.stringify(json.longUrl)
                })
            }).catch((error) => {
                console.error(error);
            });
        event.preventDefault();
    }

     

    render() {
        return (
            <div className="container">
                <div className="col-xs-8 col-xs-offset-2 jumbotron text-center">
                    <h1>URL Shortener</h1>
                    <p>Long URL to short URL</p>
                    <form onSubmit={this.handleLongURLSubmit}>
                        <label>
                            <input type="text" value={this.state.longURLValue} onChange={this.handleLongURLChange} />
                        </label>
                        <input type="submit" className="btn-primary" value="Submit" />
                    </form>
                    <p>{this.state.shortURLResult} </p>
                    <p>Short URL to long URL</p>
                    <form onSubmit={this.handleShortURLSubmit}>
                        <label>
                            <input type="text" value={this.state.shortURLValue} onChange={this.handleShortURLChange} />
                        </label>
                        <input type="submit" className="btn-primary" value="Submit" />
                    </form>
                    <p>{this.state.longURLResult} </p>
                </div>
            </div>
        )
    }
}

ReactDOM.render(<App />, document.getElementById('app'));
